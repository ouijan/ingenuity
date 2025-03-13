package main

import (
	"fmt"

	"github.com/jakecoffman/cp/v2"
	"github.com/ouijan/ingenuity/pkg/engine"
	"github.com/ouijan/ingenuity/sandbox/src/pong"
)

func main() {
	// engine.Scene.SetNext(NewDemoScene())
	engine.Scene.SetNext(pong.NewPongScene())
	engine.Run()

	// space := cp.NewSpace()
	//
	// bb := cp.BB{L: -5, B: 5, R: 5, T: -5}
	//
	// body := cp.NewBody(1, cp.INFINITY)
	// body.SetType(cp.BODY_STATIC)
	// body.SetType(cp.BODY_DYNAMIC)
	// body.SetPosition(cp.Vector{X: 50, Y: 50})
	// body.SetVelocity(-1, 0)
	// space.AddBody(body)
	//
	// shape := cp.NewBox2(body, bb, 0)
	// shape.SetElasticity(1)
	// shape.SetFriction(1)
	// space.AddShape(shape)
	//
	// dt := 1.0
	// debug(space, -1)
	// for i := 0; i < 5; i++ {
	// 	space.Step(dt)
	// 	debug(space, i)
	// }
}

func debug(space *cp.Space, i int) {
	space.EachBody(func(body *cp.Body) {
		fmt.Printf("%v ==> Body: %v, Pos: %v, Vel: %v\n", i, body, body.Position(), body.Velocity())
	})
}
