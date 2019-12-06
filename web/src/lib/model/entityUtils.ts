
export interface IName {
    name: string;
}

export function compareByName(v1: IName, v2: IName): number {
    if (v1 === undefined) return v2 === undefined ? 0 : -1;
    return v2 === undefined ? 1 : v1.name.toLocaleLowerCase().localeCompare(v2.name.toLocaleLowerCase())
}

export function sortByName<T extends IName>(itemMap: {[id: string]: T}): T[] {
    return Object.values(itemMap).sort(compareByName);
}

export interface IStringId {
    id: string;
}

export function indexById<T extends IStringId>(items: T[]): {[id: string]: T} {
    return items.reduce((byId, company) => ({...byId, [company.id]: company}), {});
}
