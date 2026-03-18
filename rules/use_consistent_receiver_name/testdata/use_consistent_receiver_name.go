package fixtures

type Server struct{}

func (s *Server) Start() {}

func (s *Server) Run() {}

func (srv *Server) Stop() {} // MATCH /receiver name srv should be consistent with previous receiver name s for Server/

type Client struct{}

func (c *Client) Connect() {}

func (c *Client) Disconnect() {}

func (cl *Client) Send() {} // MATCH /receiver name cl should be consistent with previous receiver name c for Client/

type Handler struct{}

func (_ *Handler) Handle() {}

func (_ *Handler) Process() {}

type Worker struct{}

func (w Worker) Work() {}

func (w Worker) Rest() {}

type Logger struct{}

func (l *Logger) Info() {}

func (l *Logger) Debug() {}

func (l *Logger) Warn() {}
