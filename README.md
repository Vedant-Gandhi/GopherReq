# Http-1.1-Lib

This project aims for more personal growth to learn how the HTTP protocol works under the hood, understand the complexities of the protocol and gain more insights in a step by step way. I have started this as a side project and am enthusiastic for suggestions and improvements.

To start with the development - 
1. Clone the repo to local machine.
2. Just run the command `go run main.go`.
3. Currently the socket is hardcoded at "localhost:8811" but soon will shift to env/config.
4. Things are messy now so you are welcome to refactor the code to make it more readable.
5. Feel free to open issue for any suggestion,doubt or criticism.


#### I have added support for env files. Following is a list of supported env variables - 
1. **Domain** - Set the HTTP_DOMAIN key to set your custom domain.