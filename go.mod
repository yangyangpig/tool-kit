module gather/toolkitcl

go 1.12

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190308221718-c2843e01d9a2
	golang.org/x/exp => github.com/golang/exp v0.0.0-20190121172915-509febef88a4
	golang.org/x/lint => github.com/golang/lint v0.0.0-20181026193005-c67002cb31c3
	golang.org/x/net => github.com/golang/net v0.0.0-20190320064053-1272bf9dcd53
	golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20180821212333-d2e6202438be
	golang.org/x/sync => github.com/golang/sync v0.0.0-20181108010431-42b317875d0f
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190219203350-90b0e4468f99 // indirect
	golang.org/x/text => github.com/golang/text v0.3.0
	golang.org/x/tools => github.com/golang/tools v0.0.0-20180221164845-07fd8470d635
	google.golang.org/genproto => github.com/google/go-genproto v0.0.0-20190307195333-5fe7a883aa19
	google.golang.org/grpc => github.com/grpc/grpc-go v1.33.1
	gopkg.in/yaml.v2 => github.com/go-yaml/yaml v0.0.0-20181115110504-51d6538a90f8
)

require (
	github.com/0x5010/grpcp v0.0.0-20180912032145-6d4772332891
	github.com/golang/protobuf v1.4.3
	github.com/jhump/protoreflect v1.8.1
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/panjf2000/ants/v2 v2.3.1
	github.com/panjf2000/gnet v1.0.0
	github.com/stretchr/testify v1.5.1
	google.golang.org/grpc v1.36.0
	google.golang.org/grpc/examples v0.0.0-20210513010352-dc77d7ffe311
	google.golang.org/protobuf v1.25.1-0.20200805231151-a709e31e5d12
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect

)
