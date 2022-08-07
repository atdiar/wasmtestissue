package main

import (
	"strconv"
	"sync"
	"syscall/js"
    "log"
)

// GOOS=js GOARCH=wasm go build -o  server/assets/app.wasm

var mutex = &sync.Mutex{}

// Mutable is a variable is protected for access by multiple goroutines 
// (we may want to access it asynchronously)
// It is modified in response to user actions in the browser.
// As such, it models the UI tree that would get modified by user initiated events.
// Thus, different event handlers need to take the lock before modifying it
var Mutable int 

var evtcounter = js.Global().Get("document").Call("createElement", "div")

func AddEventListener(event string, target js.Value){
	cb:= js.FuncOf(func(this js.Value, args []js.Value)any{
        // RETRIEVING EVENT INFO
        evt := args[0]
		typ := evt.Get("type").String()
        log.Print("== HANDLING " + typ + " EVENT ==\n")
        
        
        // EVENT HANDLING
        log.Print("LOCKING\n")
		b:= mutex.TryLock() // locking here as we need to access Mutable
        if !b{
            log.Print("ERROR: UNLOCKING phase seems to have been skipped in previous handler. Wrapped Functions are prempted, interleaved?")
            mutex.Lock() // just to cause the deadlock
        }
        log.Print("LOCKED\n")
		Mutable++
		evtcounter.Set("textContent", "# of event triggers: " + strconv.Itoa(Mutable))
        
		evtcounter.Call("focus")

        log.Print("UNLOCKING\n")
		mutex.Unlock() // Unlocking needs to happen before another callback starts.
        log.Print("UNLOCKED ... ready for locking\n")
		return nil
	})
	target.Call("addEventListener", event, cb)
}

func main(){
	
	body:= js.Global().Get("document").Get("body")
	button:= js.Global().Get("document").Call("createElement", "btn")
	button.Set("textContent","click me")
    
    evtcounter.Call("setAttribute", "tabindex","-1") // makes it focusable programmatically


	body.Call("appendChild",evtcounter)
	body.Call("appendChild",button)

	AddEventListener("click",button)
	AddEventListener("focusin",body)


	c := make(chan struct{}, 0)
	<-c
}

