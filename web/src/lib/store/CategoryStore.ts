import agent from 'superagent';
import {CategoryModel, ICategory} from '../model/CategoryModel';
import {compareBy, sortValues, addToMap} from '../model/entityUtils';
import {IMessageStore} from './MessageStore';
import {computed, flow, makeObservable, ObservableMap} from 'mobx';

const query = '{categories {id code description amountType parentId security income version transactionCount}}';

interface ICategoryResponse {
    body: {data: {categories: ICategory[]}};
}

const loadingCategories = 'Loading categories...';

export default class CategoryStore {
    private loading: boolean = false;
    private categoriesById = new ObservableMap<string, CategoryModel>();
    private messageStore: IMessageStore;

    constructor(messageStore: IMessageStore) {
        makeObservable(this);
        this.messageStore = messageStore;
    }

    @computed
    get categories(): CategoryModel[] {
        return sortValues(this.categoriesById, compareBy(category => category.displayName));
    }

    getCategory(id: string | number): CategoryModel {
        return this.categoriesById.get('' + id) || {} as CategoryModel;
    }

    loadCategories(): void {
        if (!this.loading && Object.keys(this.categories).length === 0) {
            this.messageStore.addProgressMessage(loadingCategories);
            this._loadCategories();
        }
    }

    private _loadCategories = flow(function* () {
        this.loading = true;
        try {
            const {body: {data}}: ICategoryResponse = yield agent.post('/finances/api/v1/graphql').send({query});
            addToMap(this.categoriesById, data.categories.map(category => new CategoryModel(category, this)));
        } catch (err) {
            console.error('error gettting categories', err); // TODO show toast
        } finally {
            this.loading = false;
            this.messageStore.removeProgressMessage(loadingCategories);
        }
    });
}