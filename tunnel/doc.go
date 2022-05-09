// Package tunnel
/*
	Client_TX:
		network_interface.read(mtu) --> process(payload) --> pack --> tunnel.write(packet)
	Client_RX:
		tunnel.read(packet) --> unpack --> process(payload) --> network_interface.write(payload)
	Server_TX:
		network_interface.read(mtu) --> identify --> route(tunn) --> pack --> tunn.write(packet)
	Server_RX:
		tunn.read(packet) --> unpack(payload) --> network_interface.write(payload)

*/
package tunnel
