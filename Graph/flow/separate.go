package flow

const fakeLevel = ^uint(0)

func (pk *data) flushBack() {
	for i := 0; i < len(pk.origin); i++ {
		origin, shadow, reflux := pk.origin[i], pk.shadow[i], pk.reflux[i]
		if len(shadow) != 0 {
			fillBackVec(origin, shadow)
			pk.shadow[i] = shadow[:0]
		}
		if len(reflux) != 0 {
			sort(reflux)
			origin = merge(origin, reflux)
		}
		pk.origin[i] = compact(origin) //去重去零
	}
}
func merge(base, part []Path) []Path {
	a, b := len(base)-1, len(part)-1
	base = append(base, part...)
	for c := len(base) - 1; b >= 0; c-- {
		if a < 0 {
			for ; b >= 0; b-- {
				base[b] = part[b]
			}
			break
		}
		if base[a].Next > part[b].Next {
			base[c] = base[a]
			a--
		} else {
			base[c] = part[b]
			b--
		}
	}
	return base
}
func compact(list []Path) []Path {
	size := len(list)
	if size == 0 {
		return list
	}
	last := 0
	for i := 1; i < size; i++ {
		if list[i].Next == list[last].Next {
			list[last].Weight += list[i].Weight
		} else {
			if list[last].Weight != 0 {
				last++
			}
			list[last] = list[i]
		}
	}
	if list[last].Weight != 0 {
		last++
	}
	return list[:last]
}

//宽度优先遍历
func (pk *data) markLevel() bool {
	for i := 0; i < len(pk.memo); i++ {
		pk.memo[i] = fakeLevel
	}

	pk.memo[pk.start] = 0
	pk.queue.push(pk.start)

	for !pk.queue.isEmpty() {
		cur := pk.queue.pop()
		for _, path := range pk.origin[cur] {
			if pk.memo[path.Next] != fakeLevel {
				continue
			}
			pk.memo[path.Next] = pk.memo[cur] + 1
			if path.Next == pk.end {
				return true
			}
			pk.queue.push(path.Next)
		}
	}
	return false
}

//筛分层次，生成分层残图，复杂度为O(E)。
func (pk *data) separate() bool {
	pk.queue.clear()
	if !pk.markLevel() {
		return false
	}
	for { //队列pop出的点并没有实际删除，可回溯遍历所有访问过的点
		cur, err := pk.queue.traceBack()
		if err != nil {
			break
		}
		paths := pk.origin[cur]
		for i := 0; i < len(paths); i++ {
			next := paths[i].Next
			if pk.memo[next] == pk.memo[cur]+1 {
				pk.shadow[cur] = append(pk.shadow[cur], paths[i])
				paths[i].Weight = 0
			}
		}
		if len(pk.shadow[cur]) == 0 {
			pk.memo[cur] = fakeLevel
		}
	}
	return true
}
