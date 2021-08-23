package model

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeometryMarshalJSONMultiPoint(t *testing.T) {
	// Given
	vals := make([]interface{}, 0)
	vals = append(vals, []interface{}{float64(1), float64(2)})
	vals = append(vals, []interface{}{float64(3), float64(4)})
	expected := "multipoint\":[[1,2],[3,4]]"

	g := Geometry{Type: GeometryMultiPoint, Coordinates: vals}
	blob, err := g.MarshalJSON()

	// When
	require.NoError(nil, err)
	assert.Equal(t, true, bytes.Contains(blob, []byte(expected)))
}

func TestGeometryMarshalJSONPoint(t *testing.T) {
	// Given
	vals := make([]interface{}, 0)
	vals = append(vals, float64(1))
	vals = append(vals, float64(2))
	expected := "point\":[1,2]"

	g := Geometry{Type: GeometryPoint, Coordinates: vals}
	blob, err := g.MarshalJSON()

	// When
	require.NoError(nil, err)
	assert.Equal(t, true, bytes.Contains(blob, []byte(expected)))
}

/*
func TestUnmarshalGeometryMultiPoint(t *testing.T) {
	rawJSON := `{"type": "MultiPoint", "coordinates": [[102.0, 0.5],[102.0, 0.5]]}`

	g := Geometry{}
	err := g.UnmarshalJSON([]byte(rawJSON))

	if err != nil {
		t.Fatalf("should unmarshal geometry without issue, err %v", err)
	}

	if g.Type != "MultiPoint" {
		t.Errorf("incorrect type, got %v", g.Type)
	}

	if len(g.MultiPoint) != 2 {
		t.Errorf("should have 2 coordinate elements but got %d", len(g.MultiPoint))
	}
}

func TestUnmarshalGeometryPoint(t *testing.T) {
	rawJSON := `{"type": "Point", "coordinates": [102.0, 0.5]}`

	g := Geometry{}
	err := g.UnmarshalJSON([]byte(rawJSON))

	if err != nil {
		t.Fatalf("should unmarshal geometry without issue, err %v", err)
	}

	if g.Type != "Point" {
		t.Errorf("incorrect type, got %v", g.Type)
	}

	if len(g.Point) != 2 {
		t.Errorf("should have 2 coordinate elements but got %d", len(g.Point))
	}
}

func TestUnmarshalGeometryMultiPolygon(t *testing.T) {
	rawJSON := `{"type": "MultiPolygon", "coordinates": [[[[1,2],[3,4]],[[5,6],[7,8]]],[[[8,7],[6,5]],[[4,3],[2,1]]]]}`

	g := &Geometry{}
	err := g.UnmarshalJSON([]byte(rawJSON))

	if err != nil {
		t.Fatalf("should unmarshal geometry without issue, err %v", err)
	}

	if g.Type != "MultiPolygon" {
		t.Errorf("incorrect type, got %v", g.Type)
	}

	if len(g.MultiPolygon) != 2 {
		t.Errorf("should have 2 polygons but got %d", len(g.MultiPolygon))
	}
}
*/
