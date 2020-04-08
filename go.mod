module github.com/atomix/cache-storage

go 1.12

require (
	github.com/atomix/api v0.0.0-20200211005812-591fe8b07ea8
	github.com/atomix/go-framework v0.0.0-20200326224656-5401b72ffe96
	github.com/atomix/go-local v0.0.0-20200326224828-5a1943478559
	github.com/atomix/kubernetes-controller v0.0.0-20200406050343-092c9b9afb2d
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/sirupsen/logrus v1.4.2
	golang.org/x/net v0.0.0-20191112182307-2180aed22343 // indirect
	golang.org/x/sys v0.0.0-20191113165036-4c7a9d0fe056 // indirect
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v0.17.2
	sigs.k8s.io/controller-runtime v0.5.2
)

replace github.com/atomix/kubernetes-controller => ../kubernetes-controller
