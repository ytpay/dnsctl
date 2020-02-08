package gdns

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/clientv3"
)

func client() *clientv3.Client {
	etcdCA, err := ioutil.ReadFile(viper.GetString("etcd.ca"))
	if err != nil {
		logrus.Fatal(err)
	}

	etcdClientCert, err := tls.LoadX509KeyPair(viper.GetString("etcd.cert"), viper.GetString("etcd.key"))
	if err != nil {
		logrus.Fatal(err)
	}

	rootCertPool := x509.NewCertPool()
	rootCertPool.AppendCertsFromPEM(etcdCA)

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   viper.GetStringSlice("etcd.endpoints"),
		DialTimeout: 5 * time.Second,
		TLS: &tls.Config{
			RootCAs:      rootCertPool,
			Certificates: []tls.Certificate{etcdClientCert},
		},
	})
	if err != nil {
		logrus.Fatal(err)
	}
	return cli
}

func PutHosts(hosts string) {
	cli := client()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key := viper.GetString("dnskey")
	if key == "" {
		logrus.Fatal("etcd gdns key set")
	}

	_, err := cli.Put(ctx, key, hosts)
	if err != nil {
		logrus.Fatal(err)
	} else {
		logrus.Info("update dns records success")
	}
}

func GetHosts() string {
	cli := client()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key := viper.GetString("dnskey")
	if key == "" {
		logrus.Fatal("etcd gdns key set")
	}

	resp, err := cli.Get(ctx, key)
	if err != nil {
		logrus.Fatal(err)
	}

	if len(resp.Kvs) != 1 {
		logrus.Error("failed to get etcd response")
		return ""
	}

	return string(resp.Kvs[0].Value)
}
