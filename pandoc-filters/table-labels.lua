function Table(element)
  element.colspecs = element.colspecs:map(function (colspec)
        local align = colspec[1]
        local width = nil -- remove width
        return {align, width}
    end)

  local heading = {}
  for _, row in pairs(element.head.rows) do
    for i, cell in ipairs(row.cells) do
      heading[i] = pandoc.utils.stringify(cell.content)
    end
  end

  for _, row in pairs(element.bodies) do
    for _, body in pairs(row.body) do
      for i, cell in ipairs(body.cells) do
        cell.attributes['label']=heading[i]
      end
    end
  end

  return element
end
