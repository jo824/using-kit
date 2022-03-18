# Using go-kit (& gizmo)

##What?
2 examples of building services in Go using Go-Kit. 

The first implementation is using raw go kit (basic mutation of the examples kit provides). and the other is leverage Gizmo
an open source library from NYT. I believe the gizmo library has it's perks out of the box, but can feel a little bloated especially for this use case.

##Why?
As we explore go-kit I believe its nice to juxtapose raw kit against gizmo. Interesting to see the lift/differences as 
each service evolves, and discussed in an accompanying blog post. Allowing someone new to Go to more easily 
reason about style/implementation choices.


###TODO
CLI/Client to easily test EPs as it becomes more complex

More about project structure.


##Running it
Project currently using Go 1.17. 
* Navigate to one of the server impl directories (gizmo or raw-kit).
* go install to install any required dependencies. (may need to edit some env variables 
 if you aren't operatoring in your go root dir. this is no longer required, but still my preference)
 the deps will end up in your `$GO_ROOT`. 
* Next head to the respective server directory & `run go build .` followed by `./server`
 and you should see it running on `localhost:8080` or `:8833`  depending on which you launched

##Notes/Worth a mention

This is  not a production level service. There are some obvious design flaws for the sake of comparing the 2.
Such as the shared code and how the map is accessed, scoping of functions (package vs global), and function/var 
definitions not being consistent, or as well thought through as they could be, in relation to them being pass by reference 
or by value.

