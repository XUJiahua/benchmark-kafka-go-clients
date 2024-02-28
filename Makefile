# install v1 ginkgo
# go install github.com/onsi/ginkgo/ginkgo@latest
# benchmark producers
bench:
	ginkgo -focus=producer ./... -- -test.bench=. -library=confluent -num=10000 -size=1000 -brokers=10.30.11.112:9092 -topic=bench
	ginkgo -focus=producer ./... -- -test.bench=. -library=sarama -num=10000 -size=1000 -brokers=10.30.11.112:9092 -topic=bench
	#ginkgo -focus=producer ./... -- -test.bench=. -library=kafkago -num=10000 -size=1000 -brokers=10.30.11.112:9092 -topic=bench

#  confluent producing 10000 messages of 1000 bytes size:
#    Fastest Time: 2.639s
#    Slowest Time: 7.734s
#    Average Time: 4.522s ± 2.283s

#  sarama producing 10000 messages of 1000 bytes size:
#    Fastest Time: 2.872s
#    Slowest Time: 3.345s
#    Average Time: 3.055s ± 0.207s

# kafka-go is extreme slow!!!
# async wins

bench-sync:
	ginkgo -focus=producer ./... -- -test.bench=. -library=sarama-sync -num=100 -size=1000 -brokers=10.30.11.112:9092 -topic=bench
	ginkgo -focus=producer ./... -- -test.bench=. -library=kafkago -num=100 -size=1000 -brokers=10.30.11.112:9092 -topic=bench

#  sarama-sync producing 100 messages of 1000 bytes size:
#    Fastest Time: 10.570s
#    Slowest Time: 10.785s
#    Average Time: 10.679s ± 0.088s

#  kafkago producing 100 messages of 1000 bytes size:
#    Fastest Time: 11.481s
#    Slowest Time: 11.507s
#    Average Time: 11.495s ± 0.011s

# batch gain performance

#  sarama-sync producing 100 messages of 100000 bytes size:
#    Fastest Time: 15.164s
#    Slowest Time: 18.286s
#    Average Time: 16.381s ± 1.364s

#  kafkago producing 100 messages of 100000 bytes size:
#    Fastest Time: 17.468s
#    Slowest Time: 19.343s
#    Average Time: 18.242s ± 0.800s

bench-batch:
	ginkgo -focus=producer ./... -- -test.bench=. -library=kafkago-batch -num=10000 -size=1000 -brokers=10.30.11.112:9092 -topic=bench

#  kafkago-batch producing 10000 messages of 1000 bytes size:
#    Fastest Time: 5.331s
#    Slowest Time: 5.810s
#    Average Time: 5.583s ± 0.197s

# produce in batch is good enough
