import agent from 'superagent';
import {CategoryModel, ICategory} from '../model/CategoryModel';
import {indexById, compareBy} from '../model/entityUtils';
import {IMessageStore} from './MessageStore';
import {computed, flow, observable} from 'mobx';

const query = '{categories {id code description amountType parentId security income version}}';

interface ICategoryResponse {
    body: {data: {categories: ICategory[]}};
}

const loadingCategories = 'Loading categories...';

export default class CategoryStore implements ICategoryStore {
    private loading: boolean = false;
    @observable
    private categoriesById: {[id: string]: CategoryModel} = {};
    private messageStore: IMessageStore;

    constructor(messageStore: IMessageStore) {
        this.messageStore = messageStore;
    }

    @computed
    get categories(): CategoryModel[] {
        return Object.values(this.categoriesById).sort(compareBy(category => category.code));
    }

    getCategory(id: string | number): CategoryModel {
        return this.categoriesById['' + id] || {} as CategoryModel;
    }

    loadCategories(): void {
        if (!this.loading && Object.keys(this.categories).length === 0) {
            this.messageStore.addProgressMessage(loadingCategories);
            this._loadCategories();
        }
    }

    private _loadCategories = flow(function*() {
        this.loading = true;
        try {
            const {body: {data}}: ICategoryResponse = yield agent.post('/finances/api/v1/graphql').send({query});
            this.categoriesById = indexById(data.categories.map(category => new CategoryModel(category, this)));
        } catch (err) {
            console.error('error gettting categories', err); // TODO show toast
        } finally {
            this.loading = false;
            this.messageStore.removeProgressMessage(loadingCategories);
        }
    });
}