package etcdhosts

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"sort"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/clientv3"
)

type vHosts struct {
	Version  int64
	Revision int64
	Hosts    string
}

type vHostsList []vHosts

func (v vHostsList) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v vHostsList) Len() int           { return len(v) }
func (v vHostsList) Less(i, j int) bool { return v[i].Version > v[j].Version }

func client() *clientv3.Client {

	home, err := homedir.Dir()
	if err != nil {
		logrus.Fatal(err)
	}

	ca := viper.GetString("etcd.ca")
	cert := viper.GetString("etcd.cert")
	key := viper.GetString("etcd.key")

	if strings.HasPrefix(ca, "~") {
		ca = strings.Replace(ca, "~", home, 1)
	}
	if strings.HasPrefix(cert, "~") {
		cert = strings.Replace(cert, "~", home, 1)
	}
	if strings.HasPrefix(key, "~") {
		key = strings.Replace(key, "~", home, 1)
	}

	etcdCA, err := ioutil.ReadFile(ca)
	if err != nil {
		logrus.Fatal(err)
	}

	etcdClientCert, err := tls.LoadX509KeyPair(cert, key)
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

func putHosts(hosts string) {
	cli := client()
	defer func() { _ = cli.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key := viper.GetString("dnskey")
	if key == "" {
		logrus.Fatal("etcd etcdhosts key set")
	}

	_, err := cli.Put(ctx, key, hosts)
	if err != nil {
		logrus.Fatal(err)
	} else {
		logrus.Info("update dns records success")
	}
}

func getHosts(revision int64) (string, error) {
	cli := client()
	defer func() { _ = cli.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key := viper.GetString("dnskey")
	if key == "" {
		logrus.Warn("etcd etcdhosts key not set")
		key = "/etcdhosts"
	}

	var resp *clientv3.GetResponse
	var err error
	if revision > -1 {
		resp, err = cli.Get(ctx, key, clientv3.WithRev(revision))
	} else {
		resp, err = cli.Get(ctx, key)
	}

	if err != nil {
		return "", err
	}

	if len(resp.Kvs) == 0 {
		return "", errors.New("etcd Hosts not exist")
	}

	if len(resp.Kvs) > 1 {
		return "", errors.New("too many etcd Hosts")
	}

	return string(resp.Kvs[0].Value), nil
}

func getHostsHistory() vHostsList {
	cli := client()
	defer func() { _ = cli.Close() }()

	key := viper.GetString("dnskey")
	if key == "" {
		key = "/etcdhosts"
		logrus.Warn("etcd etcdhosts key not set, use default key [/etcdhosts]")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	getResp, err := cli.Get(ctx, key)
	if err != nil {
		logrus.Fatal(err)
	}
	if len(getResp.Kvs) < 1 {
		logrus.Fatal("failed to get etcd data: kv not found")
	}

	vl := vHostsList{}
	for i := getResp.Header.Revision; i > 0; i-- {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		resp, err := cli.Get(ctx, key, clientv3.WithRev(i))
		if err != nil {
			break
		}
		vl = append(vl, vHosts{
			Version:  resp.Kvs[0].Version,
			Revision: i,
			Hosts:    string(resp.Kvs[0].Value),
		})
		cancel()
	}
	sort.Sort(vl)
	return vl
}
