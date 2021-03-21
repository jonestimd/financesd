import ListModel from './ListModel';
import {testAction} from 'src/test/mobxUtils';
import {createDiv} from 'src/test/htmlUtils';
import SelectionModel from './SelectionModel';

interface ITestRow {
    id: string;
    name: string;
    isValid: boolean;
    isChanged: boolean;
}

const newTestRow = (id: string) => ({id, name: `row ${id}`, isValid: true, isChanged: false});

const newModel = (...ids: string[]) => new ListModel({columns: 2}, ids.map(newTestRow));

describe('ListModel', () => {
    describe('constructor', () => {
        it('sets rows', () => {
            const items = [newTestRow('1')];

            const model = new ListModel<ITestRow>({columns: 2}, items);

            expect(model.rows).toEqual(items.length);
            expect(model.items).toEqual(items);
        });
    });
    describe('get changes', () => {
        it('excludes pending adds and deletes', () => {
            const model = newModel('1', '2', '3');
            model.cell.row = 1;
            model.delete();
            model.items[1].isChanged = true;
            model.items[2].isChanged = true;
            model.add(() => newTestRow('-1'));

            expect(model.changes).toEqual(model.items.slice(2, 3));
        });
    });
    describe('delete', () => {
        it('adds current row to pending deletes', () => {
            const model = newModel('1', '2', '3');
            model.cell.row = 1;

            testAction(() => model.pendingDeletes.length, () => model.delete());

            expect(model.pendingDeletes).toContain(model.items[1]);
            expect(model.items).toHaveLength(3);
            expect(model.rows).toEqual(3);
        });
        it('removes pending add', () => {
            const model = newModel('1', '2', '3');
            model.add(() => newTestRow('-1'));
            model.add(() => newTestRow('-2'));
            model.cell.row = 3;

            testAction(() => model.pendingAdds.length, () => model.delete());

            expect(model.pendingDeletes.length).toEqual(0);
            expect(model.pendingAdds.length).toEqual(1);
            expect(model.items).toHaveLength(4);
            expect(model.rows).toEqual(4);
            expect(model.cell.row).toEqual(3);
        });
        it('adjusts selected row', () => {
            const model = newModel('1', '2', '3');
            model.add(() => newTestRow('-1'));
            model.cell.row = 3;

            testAction(() => model.pendingAdds.length, () => model.delete());

            expect(model.cell.row).toEqual(2);
            expect(model.isChanged).toBe(false);
        });
        it('returns focus to container', () => {
            const model = newModel('1', '2', '3');
            model.cell.row = 1;
            model.container = createDiv();
            jest.spyOn(model.container, 'focus');

            testAction(() => model.pendingDeletes.length, () => model.delete());

            expect(model.container.focus).toBeCalledTimes(1);
        });
    });
    describe('undelete', () => {
        it('cancels a pending delete', () => {
            const model = newModel('1', '2', '3');
            model.cell.row = 1;
            model.delete();

            testAction(() => model.pendingDeletes.length, () => model.undelete());

            expect(model.pendingDeletes.length).toEqual(0);
            expect(model.isChanged).toBe(false);
        });
    });
    describe('isDelete', () => {
        it('returns true if current item is a pending delete', () => {
            const model = newModel('1', '2', '3');
            model.cell.row = 1;
            model.delete();

            expect(model.isDelete()).toBe(true);
        });
        it('returns true if item is a pending delete', () => {
            const model = newModel('1', '2', '3');
            model.cell.row = 1;
            model.delete();

            expect(model.isDelete(model.items[0])).toBe(false);
            expect(model.isDelete(model.items[1])).toBe(true);
            expect(model.isDelete(model.items[2])).toBe(false);
        });
    });
    describe('selected', () => {
        it('returns the selected item', () => {
            const model = newModel('1', '2', '3');
            model.cell.row = 1;

            expect(model.selected).toBe(model.items[1]);
        });
    });
    describe('items', () => {
        it('saved items and pending adds', () => {
            const model = newModel('1', '2', '3');
            model.add(() => newTestRow('-1'));

            expect(model.items).toEqual(model['_items'].concat(model.pendingAdds));
        });
    });
    describe('isChanged', () => {
        it('returns false if no changes, adds or deletes', () => {
            expect(newModel('1', '2', '3').isChanged).toBe(false);
        });
        it('returns true if an item is changed', () => {
            const model = newModel('1', '2', '3');
            model.items[1].isChanged = true;

            expect(model.isChanged).toBe(true);
        });
        it('returns true if an item is added', () => {
            const model = newModel('1', '2', '3');
            model.add(() => newTestRow('-1'));

            expect(model.isChanged).toBe(true);
        });
        it('returns true if an item is deleted', () => {
            const model = newModel('1', '2', '3');
            model.delete();

            expect(model.isChanged).toBe(true);
        });
    });
    describe('isValid', () => {
        it('returns true if all items are valid', () => {
            const model = newModel('1', '2', '3');
            model.add(() => newTestRow('-1'));

            expect(model.isValid).toBe(true);
        });
        it('returns false if an item is not valid', () => {
            const model = newModel('1', '2', '3');
            model.items[1].isValid = false;

            expect(model.isValid).toBe(false);
        });
        it('returns false if a pending add is not valid', () => {
            const model = newModel('1', '2', '3');
            model.add(() => newTestRow('-1'));
            model.pendingAdds[0].isValid = false;

            expect(model.isValid).toBe(false);
        });
    });
    describe('isEditable', () => {
        it('returns false for pending delete', () => {
            const model = newModel('1', '2', '3');
            model.cell.row = 1;
            model.delete();

            expect(model.isEditable(model.items[0])).toBe(true);
            expect(model.isEditable(model.items[1])).toBe(false);
            expect(model.isEditable(model.items[2])).toBe(true);
        });
    });
    describe('rowClass', () => {
        it('calls SelectionModel.rowClass', () => {
            const model = newModel('1', '2', '3');
            jest.spyOn(SelectionModel.prototype, 'rowClass').mockReturnValue('selected');

            expect(model.rowClass(0)).toEqual('selected');
        });
        it('appends deleted to pending delete', () => {
            const model = newModel('1', '2', '3');
            model.delete();
            jest.spyOn(SelectionModel.prototype, 'rowClass').mockReturnValue('selected');

            expect(model.rowClass(0)).toEqual('selected deleted');
            expect(model.rowClass(1)).toEqual('selected');
            expect(model.rowClass(2)).toEqual('selected');
        });
    });
    describe('commit', () => {
        it('clears pending changes and replaces items', () => {
            const model = newModel('1', '2', '3');
            model.delete();
            model.add(() => newTestRow('-1'));
            const result = model.items.slice(1, 3).concat(newTestRow('4'));

            testAction(() => model.pendingAdds.length, () => model.commit(result));

            expect(model.items).toEqual(result);
            expect(model.isChanged).toBe(false);
        });
    });
});
