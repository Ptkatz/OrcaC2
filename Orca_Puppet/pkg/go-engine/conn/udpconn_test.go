package conn

import (
	"fmt"
	"Orca_Puppet/pkg/go-engine/loggo"
	"strconv"
	"testing"
	"time"
)

func Test000UDP(t *testing.T) {
	c, err := NewConn("udp")
	if err != nil {
		fmt.Println(err)
		return
	}

	cc, err := c.Listen(":58086")
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		cc.Accept()
		fmt.Println("accept done")
	}()

	time.Sleep(time.Second)

	cc.Close()

	time.Sleep(time.Second)
}

func Test0002UDP(t *testing.T) {
	c, err := NewConn("udp")
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		conn, err := c.Dial("9.9.9.9:58086")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(conn.Info())
		}

	}()

	time.Sleep(time.Second)

	c.Close()

	time.Sleep(time.Second)
}

func Test0003UDP(t *testing.T) {
	c, err := NewConn("udp")
	if err != nil {
		fmt.Println(err)
		return
	}

	cc, err := c.Listen(":58086")
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		cc.Accept()
		fmt.Println("accept done")
	}()

	ccc, err := c.Dial(":58086")
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		buf := make([]byte, 100)
		_, err := ccc.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	time.Sleep(time.Second)

	cc.Close()
	ccc.Close()

	time.Sleep(time.Second)
}

func Test0004UDP(t *testing.T) {
	c, err := NewConn("udp")
	if err != nil {
		fmt.Println(err)
		return
	}

	cc, err := c.Listen(":58086")
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		cc.Accept()
		fmt.Println("accept done")
	}()

	ccc, err := c.Dial(":58086")
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		buf := make([]byte, 1000)
		for i := 0; i < 10000; i++ {
			_, err := ccc.Write(buf)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		fmt.Println("write done")
	}()

	time.Sleep(time.Second)

	cc.Close()
	ccc.Close()

	time.Sleep(time.Second)
}

func Test0005UDP(t *testing.T) {
	c, err := NewConn("udp")
	if err != nil {
		fmt.Println(err)
		return
	}

	cc, err := c.Listen(":58086")
	if err != nil {
		fmt.Println(err)
		return
	}

	exit := false

	go func() {
		cc, err := cc.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		defer cc.Close()
		fmt.Println("accept done")
		buf := make([]byte, 10)
		for !exit {
			n, err := cc.Read(buf)
			if err != nil {
				fmt.Println(err)
				fmt.Println("Read done")
				return
			}
			fmt.Println(string(buf[0:n]))
			time.Sleep(time.Millisecond * 100)
		}
	}()

	ccc, err := c.Dial(":58086")
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for i := 0; i < 10000 && !exit; i++ {
			_, err := ccc.Write([]byte("hahaha" + strconv.Itoa(i)))
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		fmt.Println("write done")
	}()

	time.Sleep(time.Second)

	cc.Close()
	ccc.Close()

	exit = true

	time.Sleep(time.Second)
}

func Test0005UDP1(t *testing.T) {
	c, err := NewConn("udp")
	if err != nil {
		fmt.Println(err)
		return
	}

	cc, err := c.Listen(":58086")
	if err != nil {
		fmt.Println(err)
		return
	}

	exit := false

	go func() {
		fmt.Println("start Accept")
		cc, err := cc.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("accept done")
		for i := 0; i < 10000 && !exit; i++ {
			_, err := cc.Write([]byte("hahaha" + strconv.Itoa(i)))
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		fmt.Println("write done")
	}()

	fmt.Println("start Dial")
	ccc, err := c.Dial(":58086")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Dial done")

	go func() {
		buf := make([]byte, 10)
		ccc.Write([]byte("hahaha"))
		ccc.Write([]byte("hahaha"))
		ccc.Write([]byte("hahaha"))
		for {
			n, err := ccc.Read(buf)
			if err != nil {
				fmt.Println(err)
				fmt.Println("Read done")
				return
			}
			fmt.Println(string(buf[0:n]))
			time.Sleep(time.Millisecond * 100)
		}
		fmt.Println("write done")
	}()

	time.Sleep(time.Second)

	cc.Close()
	ccc.Close()

	exit = true

	time.Sleep(time.Second)
}

func Test0008UDP(t *testing.T) {
	c, err := NewConn("udp")
	if err != nil {
		fmt.Println(err)
		return
	}

	cc, err := c.Listen(":58086")
	if err != nil {
		fmt.Println(err)
		return
	}

	exit := false

	go func() {
		cc, err := cc.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("accept done")
		data := make([]byte, 500)
		start := time.Now()
		speed := 0
		for !exit {
			//fmt.Println("start Write")
			_, err := cc.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
			//fmt.Println("end Write")
			speed += len(data)
			if time.Now().Sub(start) > time.Second {
				speed = speed / 1024 / 1024
				loggo.Info("write speed %v MB per second", float64(speed)/float64(time.Now().Sub(start)/time.Second))
				speed = 0
				start = time.Now()
			}
		}
		fmt.Println("write done")
	}()

	ccc, err := c.Dial(":58086")
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		fmt.Println("start client")
		ccc.Write([]byte("hahaha"))
		ccc.Write([]byte("hahaha"))
		ccc.Write([]byte("hahaha"))
		buf := make([]byte, 500)
		start := time.Now()
		speed := 0
		for !exit {
			//fmt.Println("start Read")
			n, err := ccc.Read(buf)
			//fmt.Println("start Read")
			if err != nil {
				fmt.Println(err)
				fmt.Println("Read done")
				return
			}
			speed += n
			if time.Now().Sub(start) > time.Second {
				speed = speed / 1024 / 1024
				loggo.Info("read speed %v MB per second", float64(speed)/float64(time.Now().Sub(start)/time.Second))
				speed = 0
				start = time.Now()
			}
		}
		fmt.Println("write done")
	}()

	time.Sleep(time.Second * 10)

	cc.Close()
	ccc.Close()

	exit = true

	time.Sleep(time.Second)
}
