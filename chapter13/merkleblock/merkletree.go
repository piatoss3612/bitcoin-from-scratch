package merkleblock

import (
	"chapter13/utils"
	"fmt"
	"math"
	"strings"
)

type MerkleTree struct {
	Total        int
	MaxDepth     int
	Nodes        [][][]byte
	CurrentDepth int
	CurrentIndex int
}

func NewMerkleTree(total int) *MerkleTree {
	maxDepth := int(math.Ceil(math.Log2(float64(total))))

	nodes := make([][][]byte, maxDepth+1)

	for i := 0; i <= maxDepth; i++ {
		numItems := int(math.Ceil(float64(total) / math.Pow(2, float64(maxDepth-i))))
		levelHahes := make([][]byte, numItems)
		nodes[i] = levelHahes
	}

	return &MerkleTree{
		Total:        total,
		MaxDepth:     maxDepth,
		Nodes:        nodes,
		CurrentDepth: 0,
		CurrentIndex: 0,
	}
}

func (m MerkleTree) String() string {
	builder := strings.Builder{}

	for depth, level := range m.Nodes {
		items := []string{}
		for idx, h := range level {
			var item string
			if len(h) == 0 {
				item = "None"
			} else {
				hexHash := fmt.Sprintf("%x", h)
				item = fmt.Sprintf("%s...", hexHash[:8])
			}

			if depth == m.CurrentDepth && idx == m.CurrentIndex {
				items = append(items, fmt.Sprintf("*%s*", item))
			} else {
				items = append(items, item)
			}
		}

		builder.WriteString(strings.Join(items, ","))
		builder.WriteString("\n")
	}

	return builder.String()
}

func (m *MerkleTree) Up() *MerkleTree {
	m.CurrentDepth--
	m.CurrentIndex /= 2

	return m
}

func (m *MerkleTree) Left() *MerkleTree {
	m.CurrentDepth++
	m.CurrentIndex *= 2

	return m
}

func (m *MerkleTree) Right() *MerkleTree {
	m.CurrentDepth++
	m.CurrentIndex = m.CurrentIndex*2 + 1

	return m
}

func (m *MerkleTree) Root() []byte {
	return m.Nodes[0][0]
}

func (m *MerkleTree) SetCurrentNode(hash []byte) *MerkleTree {
	m.Nodes[m.CurrentDepth][m.CurrentIndex] = hash

	return m
}

func (m *MerkleTree) GetCurrentNode() []byte {
	return m.Nodes[m.CurrentDepth][m.CurrentIndex]
}

func (m *MerkleTree) GetLeftNode() []byte {
	return m.Nodes[m.CurrentDepth+1][m.CurrentIndex*2]
}

func (m *MerkleTree) GetRightNode() []byte {
	return m.Nodes[m.CurrentDepth+1][m.CurrentIndex*2+1]
}

func (m *MerkleTree) IsLeaf() bool {
	return m.CurrentDepth == m.MaxDepth
}

func (m *MerkleTree) RightExists() bool {
	return len(m.Nodes[m.CurrentDepth+1]) > m.CurrentIndex*2+1
}

func (m *MerkleTree) PopulateTree(flagBits []byte, hashes [][]byte) error {
	for m.Root() == nil {
		if m.IsLeaf() {
			m.SetCurrentNode(hashes[0])
			hashes = hashes[1:]
			flagBits = flagBits[1:]
			m.Up()
		} else {
			leftHash := m.GetLeftNode()
			if leftHash == nil {
				bit := flagBits[0]
				flagBits = flagBits[1:]

				if bit == 0 {
					m.SetCurrentNode(hashes[0])
					hashes = hashes[1:]
					m.Up()
				} else {
					m.Left()
				}
			} else if m.RightExists() {
				rightHash := m.GetRightNode()
				if rightHash == nil {
					m.Right()
				} else {
					m.SetCurrentNode(utils.MerkleParent(leftHash, rightHash))
					m.Up()
				}
			} else {
				m.SetCurrentNode(utils.MerkleParent(leftHash, leftHash))
				m.Up()
			}
		}
	}

	if len(hashes) != 0 {
		return fmt.Errorf("hashes not all consumed %d", len(hashes))
	}

	for _, bit := range flagBits {
		if bit != 0 {
			return fmt.Errorf("flagBits not all consumed")
		}
	}

	return nil
}
