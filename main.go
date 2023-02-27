package main

func main() {
	c := New()
	c.Connect()
	c.Ping()
	c.DummyData()
	c.Disconnect()
}
