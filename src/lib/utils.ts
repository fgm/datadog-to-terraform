type Converter = (k: string, v: any) => object | string

function literalString(value: any) {
  if (typeof value === "string") {
    if (value.includes("\n")) {
      return `<<EOF\n${value}\nEOF`;
    }
    return `"${value}"`;
  } else if (Array.isArray(value)) {
    let result = "[";
    value.forEach((elem, index) => {
      result += literalString(elem);
      if (index !== value.length - 1) result += ",";
    });
    return result + "]";
  }
  return value;
}

export function assignmentString(key: string, value: any): string {
  if (value === null) return "";
  const displayValue = literalString(value);
  return `\n${key} = ${displayValue}`;
}

export function map(contents: object, converter: Converter) {
  let result = "";
  Object.entries(contents).forEach(([key, value]) => {
    result += converter(key, value);
  });
  return result;
}

export function block(name: string, contents: object, converter: Converter) {
  return `\n${name} {${map(contents, converter)}\n}`;
}

export function blockList(array: Array<object>, blockName: string, contentConverter: Converter) {
  let result = "";
  array.forEach((elem: object) => {
    result += block(blockName, elem, contentConverter);
  });
  return result;
}

export function convertFromDefinition(definitionSet: any, k: string, v: any) {
  if (typeof definitionSet[k] !== "function") throw `Can't convert key '${k}'`;
  return definitionSet[k](v);
}
