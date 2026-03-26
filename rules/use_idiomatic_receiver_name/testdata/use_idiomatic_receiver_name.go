package fixtures

type Server struct{}

func (this *Server) Start() {} // MATCH /receiver name should be a reflection of its identity; don't use generic names such as "this"/

func (self *Server) Stop() {} // MATCH /receiver name should be a reflection of its identity; don't use generic names such as "self"/

type Client struct{}

func (c *Client) Connect() {}

func (c *Client) Disconnect() {}

func (cl *Client) Send() {}

type Handler struct{}

func (_ *Handler) Handle() {} // MATCH /receiver name should not be an underscore, omit the name if it is unused/

type Worker struct{}

func (w Worker) Work() {}

func (w Worker) Rest() {}
