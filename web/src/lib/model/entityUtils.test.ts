import {addToMap, compareBy, compareByName, IName, sortValuesByName} from './entityUtils';

describe('entityUtils', () => {
    const item1 = {name: 'abc'};
    const item2 = {name: 'xyz'};

    describe('compareBy', () => {
        it('returns comparator that compares values returned by getter', () => {
            const getter = (item: IName) => item.name;

            const comparator = compareBy(getter);

            expect(comparator(item1, item2)).toBeLessThan(0);
            expect(comparator(item2, item1)).toBeGreaterThan(0);
            expect(comparator(item1, item1)).toEqual(0);
            expect(comparator(item2, item2)).toEqual(0);
        });
    });
    describe('compareByName', () => {
        it('compares name property', () => {
            expect(compareByName(item1, item2)).toBeLessThan(0);
            expect(compareByName(item2, item1)).toBeGreaterThan(0);
            expect(compareByName(item1, item1)).toEqual(0);
            expect(compareByName(item2, item2)).toEqual(0);
        });
        it('sorts undefined first', () => {
            expect(compareByName(undefined, item1)).toBeLessThan(0);
            expect(compareByName(undefined, item1)).toBeLessThan(0);
            expect(compareByName(item1, undefined)).toBeGreaterThan(0);
            expect(compareByName(item2, undefined)).toBeGreaterThan(0);
            expect(compareByName(undefined, undefined)).toEqual(0);
        });
    });
    describe('sortValuesByName', () => {
        it('sorts map values by name property', () => {
            const map = new Map<string, IName>();
            map.set('a', item2);
            map.set('b', item1);

            expect(sortValuesByName(map)).toEqual([item1, item2]);
        });
    });
    describe('addToMap', () => {
        it('adds items to map using id property', () => {
            const entry1 = {id: '1', name: 'abc'};
            const entry2 = {id: '2', name: 'xyz'};

            const map = new Map<string, typeof entry1>();
            addToMap(map, [entry1, entry2]);

            expect(map.get(entry1.id)).toBe(entry1);
            expect(map.get(entry2.id)).toBe(entry2);
        });
    });
});
