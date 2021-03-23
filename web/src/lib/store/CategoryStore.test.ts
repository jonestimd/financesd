import {newCategory, newCategoryModel} from 'src/test/categoryFactory';
import {RootStore} from './RootStore';
import * as entityUtils from '../model/entityUtils';
import * as agent from '../agent';
import {loadingCategories, query} from './CategoryStore';
import {CategoryModel} from '../model/CategoryModel';

describe('CategoryStore', () => {
    const {categoryStore, messageStore, alertStore} = new RootStore();

    beforeEach(() => {
        categoryStore['categoriesById'].clear();
    });
    describe('get categories', () => {
        it('sorts by displayName', () => {
            const displayName = 'category display name';
            const category = newCategoryModel({categoryStore});
            jest.spyOn(category, 'displayName', 'get').mockReturnValue(displayName);
            categoryStore['categoriesById'].set(category.id, category);
            const comparator = jest.fn().mockReturnValue(0);
            jest.spyOn(entityUtils, 'sortValues');
            const compareBy = jest.spyOn(entityUtils, 'compareBy').mockReturnValue(comparator);

            const categories = categoryStore.categories;

            expect(categories).toEqual([category]);
            expect(entityUtils.sortValues).toBeCalledWith(categoryStore['categoriesById'], comparator);
            expect(compareBy.mock.calls[0][0](category)).toEqual(displayName);
        });
    });
    describe('getCategory', () => {
        it('returns category for ID', () => {
            const category = newCategoryModel();
            categoryStore['categoriesById'].set(category.id, category);

            expect(categoryStore.getCategory(category.id)).toBe(category);
            expect(categoryStore.getCategory(parseInt(category.id))).toBe(category);
        });
        it('returns undefined for unknown ID', () => {
            expect(categoryStore.getCategory('-99')).toBeUndefined();
        });
    });
    describe('loadCategories', () => {
        beforeEach(() => {
            categoryStore['loading'] = false;
            jest.spyOn(messageStore, 'addProgressMessage');
            jest.spyOn(messageStore, 'removeProgressMessage');
        });
        it('loads categories if categoriesById is empty', async () => {
            const category = newCategory();
            jest.spyOn(agent, 'graphql').mockResolvedValue({data: {categories: [category]}});

            await categoryStore.loadCategories();

            expect(categoryStore['loading']).toBe(false);
            expect(messageStore.addProgressMessage).toBeCalledWith(loadingCategories);
            expect(messageStore.removeProgressMessage).toBeCalledWith(loadingCategories);
            expect(agent.graphql).toBeCalledWith(query, undefined);
            expect(categoryStore.categories).toStrictEqual([new CategoryModel(category, categoryStore)]);
        });
        it('does nothing is already loading', async () => {
            categoryStore['loading'] = true;
            jest.spyOn(agent, 'graphql').mockRejectedValue(new Error());

            await categoryStore.loadCategories();

            expect(categoryStore['loading']).toBe(true);
            expect(messageStore.addProgressMessage).not.toBeCalled();
            expect(messageStore.removeProgressMessage).not.toBeCalled();
            expect(agent.graphql).not.toBeCalled();
        });
        it('does nothing is already loaded', async () => {
            const category = newCategoryModel();
            categoryStore['categoriesById'].set(category.id, category);
            jest.spyOn(agent, 'graphql').mockRejectedValue(new Error());

            await categoryStore.loadCategories();

            expect(categoryStore['loading']).toBe(false);
            expect(messageStore.addProgressMessage).not.toBeCalled();
            expect(messageStore.removeProgressMessage).not.toBeCalled();
            expect(agent.graphql).not.toBeCalled();
        });
        it('logs error from graphql', async () => {
            const error = new Error('API error');
            jest.spyOn(agent, 'graphql').mockRejectedValue(error);
            jest.spyOn(console, 'error').mockImplementation(() => { });
            jest.spyOn(alertStore, 'addAlert').mockReturnValue();

            await categoryStore.loadCategories();

            expect(categoryStore['loading']).toBe(false);
            expect(messageStore.addProgressMessage).toBeCalledWith(loadingCategories);
            expect(messageStore.removeProgressMessage).toBeCalledWith(loadingCategories);
            expect(alertStore.addAlert).toBeCalledWith('error', 'Error loading categories');
            expect(console.error).toBeCalledWith('error from Loading categories', error);
        });
    });
});
