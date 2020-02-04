-- tab:目标表
-- len:表的长度
function IsArraySortedOrder(tab, len)
	if(len == 1) then return true end
	if(tab[len] < tab[len - 1]) then return false end
	return IsArraySortedOrder(tab, len - 1)
end

tab = {3, 2, 4, 1, 7, 5, 8, 0, 14, 11, 10}
print(IsArraySortedOrder(tab, #tab))