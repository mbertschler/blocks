package blocks

type Blocks []Block

func (Blocks) Render() Block { return nil }

type Block interface {
	Render() Block
}
