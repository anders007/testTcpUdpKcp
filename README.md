# testTcpUdpKcp
test tcp/udp/kcp rtt &amp; lostrate


depends:
go get -u github.com/golang/protobuf/proto
go get -u github.com/xtaci/kcp-go

在我的vultr东京节点上测试，KCP优势不明显，可能网络情况比较良好
KCP官方提供了一个模拟丢包和延迟的模拟器
不过我认为还是在公网真实环境下测试比较好点，有时间再弄


裸UDP公网丢包率30% - 40%，日了狗了~~~~~~
