 package socket

 //Middleware ..
 type Middleware interface {
	 process(next Context) error
 }
 