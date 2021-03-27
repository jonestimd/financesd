import {CategoryModel, ICategory} from '../model/CategoryModel';
import {compareBy, sortValues, addToMap} from '../model/entityUtils';
import {IMessageStore} from './MessageStore';
import {computed, makeObservable, ObservableMap} from 'mobx';
import AlertStore from './AlertStore';
import Loader from './Loader';

export const query = '{categories {id code description amountType parentId security income version transactionCount}}';

export const loadingCategories = 'Loading categories';

export default class CategoryStore {
    private loading = false;
    private categoriesById = new ObservableMap<number, CategoryModel>();
    private loader: Loader;

    constructor(messageStore: IMessageStore, alertStore: AlertStore) {
        makeObservable(this);
        this.loader = new Loader(messageStore, alertStore);
    }

    @computed
    get categories(): CategoryModel[] {
        return sortValues(this.categoriesById, compareBy((category) => category.displayName));
    }

    getCategory(id?: number): CategoryModel | undefined {
        return typeof id === 'number' ? this.categoriesById.get(id) : undefined;
    }

    loadCategories(): Promise<boolean> | undefined {
        if (!this.loading && Object.keys(this.categories).length === 0) {
            this.loading = true;
            return this.loader.load<{categories: ICategory[]}>(loadingCategories, {query,
                updater: ({categories}) => addToMap(this.categoriesById, categories.map((category) => new CategoryModel(category, this))),
                completer: () => this.loading = false,
            });
        }
    }
}
