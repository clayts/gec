package graphics

type UniformLocation int32

func (u UniformLocation) GL() int32 { return int32(u) }
