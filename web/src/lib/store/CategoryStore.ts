import * as agent from '../agent';
import {CategoryModel, ICategory} from '../model/CategoryModel';
import {compareBy, sortValues, addToMap} from '../model/entityUtils';
import {IMessageStore} from './MessageStore';
import {computed, flow, makeObservable, ObservableMap} from 'mobx';
import {LoadResult} from './interfaces';

export const query = '{categories {id code description amountType parentId security income version transactionCount}}';

type CategoryResponse = agent.IGraphqlResponse<{categories: ICategory[]}>;

export const loadingCategories = 'Loading categories...';

export default class CategoryStore {
    private loading = false;
    private categoriesById = new ObservableMap<string, CategoryModel>();
    private messageStore: IMessageStore;

    constructor(messageStore: IMessageStore) {
        makeObservable(this);
        this.messageStore = messageStore;
    }

    @computed
    get categories(): CategoryModel[] {
        return sortValues(this.categoriesById, compareBy((category) => category.displayName));
    }

    getCategory(id?: string | number): CategoryModel | undefined {
        return this.categoriesById.get('' + id);
    }

    loadCategories(): Promise<void> | undefined {
        if (!this.loading && Object.keys(this.categories).length === 0) {
            this.messageStore.addProgressMessage(loadingCategories);
            return this._loadCategories();
        }
    }

    private _loadCategories = flow(function* (this: CategoryStore): LoadResult<CategoryResponse> {
        this.loading = true;
        try {
            const {data} = yield agent.graphql(query);
            addToMap(this.categoriesById, data.categories.map((category) => new CategoryModel(category, this)));
        } catch (err) {
            console.error('error gettting categories', err); // TODO show toast
        } finally {
            this.loading = false;
            this.messageStore.removeProgressMessage(loadingCategories);
        }
    });
}
