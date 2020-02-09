package etcdhosts

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"sort"
	"time"

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

func getHosts() string {
	cli := client()
	defer func() { _ = cli.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key := viper.GetString("dnskey")
	if key == "" {
		logrus.Warn("etcd etcdhosts key not set")
		key = "/etcdhosts"
	}

	resp, err := cli.Get(ctx, key)
	if err != nil {
		logrus.Error(err)
		return ""
	}

	if len(resp.Kvs) == 0 {
		logrus.Error("etcd Hosts not exist")
		return ""
	}

	if len(resp.Kvs) > 1 {
		logrus.Error("too many etcd Hosts")
		return ""
	}

	return string(resp.Kvs[0].Value)
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
	ch := cli.Watch(context.Background(), key, clientv3.WithRev(getResp.Kvs[0].CreateRevision))
	for {
		select {
		case resp := <-ch:
			for _, e := range resp.Events {
				vl = append(vl, vHosts{
					Version:  e.Kv.Version,
					Revision: e.Kv.ModRevision,
					Hosts:    string(e.Kv.Value),
				})
			}
		case <-time.After(300 * time.Millisecond):
			goto done
		}
	}
done:
	sort.Sort(vl)
	return vl
}
