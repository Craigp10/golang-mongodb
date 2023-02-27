package main

func main() {
	c := New()
	c.Connect()
	c.DummyData()
	c.Disconnect()
}
