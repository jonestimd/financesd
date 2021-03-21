import {action, computed, IObservableArray, makeObservable, observable} from 'mobx';
import SelectionModel, {IOptions} from './SelectionModel';

let nextId = 1;

export interface IItem {
    readonly isValid: boolean;
    readonly isChanged: boolean;
}

export default class ListModel<T extends IItem> extends SelectionModel {
    private readonly _items: IObservableArray<T>;
    private readonly _pendingAdds = observable.array<T>();
    private readonly _pendingDeletes = observable.set<T>();

    constructor(options: Omit<IOptions, 'rows'>, items: T[]) {
        super({...options, rows: items.length});
        this._items = observable.array(items);
        makeObservable(this);
    }

    @action
    add(supplier: (id: number) => T) {
        this._pendingAdds.push(supplier(nextId++));
        this.rows++;
    }

    get pendingAdds() {
        return this._pendingAdds.slice();
    }

    get pendingDeletes() {
        return Array.from(this._pendingDeletes);
    }

    get changes() {
        return this._items.filter((c) => !this.isDelete(c)).filter((c) => c.isChanged);
    }

    @action
    delete() {
        const item = this.items[this.cell.row];
        if (this._pendingAdds.remove(item)) {
            this.rows--;
            if (this.cell.row >= this.rows) this.cell.row = this.rows - 1;
        }
        else this._pendingDeletes.add(item);
        this.container?.focus();
    }

    @action
    undelete() {
        const item = this.items[this.cell.row];
        this._pendingDeletes.delete(item);
    }

    isDelete(item: T | undefined = this.selected) {
        return this._pendingDeletes.has(item);
    }

    @computed
    get selected(): T | undefined {
        return this.items[this.cell.row];
    }

    @computed
    get items() {
        return this._items.concat(this._pendingAdds);
    }

    @computed
    get isChanged() {
        return this._items.some((i) => i.isChanged) || this._pendingDeletes.size > 0 || this._pendingAdds.length > 0;
    }

    @computed
    get isValid() {
        return !this._items.some((i) => !i.isValid) && !this._pendingAdds.some((i) => !i.isValid);
    }

    isEditable(item: T) {
        return !this._pendingDeletes.has(item);
    }

    rowClass(index: number) {
        let className = SelectionModel.prototype.rowClass.call(this, index);
        if (this._pendingDeletes.has(this.items[index])) className += ' deleted';
        return className.trim();
    }

    @action
    commit(items: T[]) {
        this._pendingAdds.clear();
        this._pendingDeletes.clear();
        this._items.replace(items);
    }
}
