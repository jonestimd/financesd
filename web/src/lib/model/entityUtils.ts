export interface IName {
    name: string;
}

export type Comparator<T> = (v1: T, v2: T) => number;

export function compareBy<T, P>(getter: (item: T) => P): Comparator<T> {
    return (v1: T, v2: T) => {
        const p1 = getter(v1);
        const p2 = getter(v2);
        if (p1 === p2) return 0;
        return p1 < p2 ? -1 : 1;
    };
}

export function compareByName(v1?: IName, v2?: IName): number {
    if (v1 === undefined) return v2 === undefined ? 0 : -1;
    return v2 === undefined ? 1 : v1.name.toLocaleLowerCase().localeCompare(v2.name.toLocaleLowerCase());
}

export function sortValuesByName<T extends IName>(itemMap: Map<string, T>): T[] {
    return sortValues<T>(itemMap, compareByName);
}

export function sortValues<T>(itemMap: Map<string, T>, comparator: Comparator<T>): T[] {
    return Array.from(itemMap.values()).sort(comparator);
}

export interface IStringId {
    id: string;
}

export function addToMap<T extends IStringId>(byId: Map<string, T>, items: T[]): void {
    items.forEach((item) => byId.set(item.id, item));
}
