export interface IName {
    name: string;
}

export function compareBy<T, P>(getter: (item: T) => P): (v1: T, v2: T) => number {
    return (v1: T, v2: T) => {
        const p1 = getter(v1);
        const p2 = getter(v2);
        if (p1 === p2) return 0;
        return p1 < p2 ? -1 : 1;
    };
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
