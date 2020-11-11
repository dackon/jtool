package schema8

import (
	"context"
	"fmt"
	"net/url"

	"github.com/dackon/jtool/jvalue"
)

// MatchInfo
type MatchInfo struct {
	Err                     error
	KeywordLocation         string
	InstanceLocation        string
	AbsoluteKeywordLocation string
}

func (mc *MatchInfo) String() string {
	return fmt.Sprintf("[Err: %s. KeywordLocation: %s. "+
		"InstanceLocation: %s. AbsoluteKeywordLocation: %s]", mc.Err,
		mc.KeywordLocation, mc.InstanceLocation, mc.AbsoluteKeywordLocation)
}

func setFailedMatchInfo(sn *schemaNode, v *jvalue.V, err error,
	ctx context.Context) {

	mc := ctx.Value(matchContextKey).(*matchContext)
	mc.mi.Err = err
	mc.mi.KeywordLocation = getKeyworadLocation(mc.path)
	mc.mi.AbsoluteKeywordLocation = getAbsoluteKeywordLocation(sn)
	mc.mi.InstanceLocation = v.GetJSONPointer()
}

func getKeyworadLocation(nodes []*matchPathNode) string {
	pns := make([]*pathNode, len(nodes))
	for i, v := range nodes {
		pns[i] = v.path
	}

	lastNode := nodes[len(nodes)-1]
	if lastNode.assertionName != "" {
		pns = append(pns, newPathNodeString(lastNode.assertionName))
	}

	return pathToFragment(pns)
}

func getAbsoluteKeywordLocation(sn *schemaNode) string {
	if len(sn.path) == 0 {
		return ""
	}

	var arr []*pathNode

	node := sn
	for node != nil {
		arr = append(arr, node.path[len(node.path)-1])
		if node.isRoot() {
			break
		}
		node = node.parent
	}

	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}

	fragment := pathToFragment(arr)
	sub, err := url.Parse(fragment)
	if err != nil {
		panic(fmt.Errorf("fragment is %s", fragment))
	}
	return node.baseURIObj.ResolveReference(sub).String()
}
