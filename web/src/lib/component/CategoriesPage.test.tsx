import React from 'react';
import {shallow} from 'enzyme';
import {RootStore} from '../store/RootStore';
import {newCategoryModel} from 'src/test/categoryFactory';
import CategoriesPage from './CategoriesPage';
import Table, {IColumn} from './table/Table';
import TopAppBar from './TopAppBar';
import {CategoryModel} from '../model/CategoryModel';
import amountType from '../i18n/amountType';
import {ObservableMap} from 'mobx';

describe('CategoriesPage', () => {
    const rootStore = new RootStore();
    const {categoryStore} = rootStore;
    const parent = newCategoryModel({categoryStore});
    const category = newCategoryModel({categoryStore, parentId: parseInt(parent.id)});

    beforeEach(() => {
        jest.spyOn(React, 'useContext').mockReturnValue(rootStore);
    });
    it('displays app bar and table of categories', () => {
        const categories = [parent, category];
        jest.spyOn(categoryStore, 'categories', 'get').mockReturnValue(categories);

        const component = shallow(<CategoriesPage />);

        expect(component.find(TopAppBar)).toHaveProp('title', 'Categories');
        expect(component.find(Table)).toHaveProp('data', categories);
    });
    describe('categories table', () => {
        beforeAll(() => {
            categoryStore['categoriesById'] = new ObservableMap([
                [parent.id, parent],
                [category.id, category],
            ]);
        });

        const columnTests: {input: CategoryModel, key: string, value: React.ReactNode, name?: string}[] = [
            {input: category, key: 'parent', value: parent.displayName},
            {input: parent, key: 'parent', value: null, name: 'displays blank for no parent'},
            {input: category, key: 'code', value: category.code},
            {input: category, key: 'description', value: category.description},
            {input: category, key: 'amountType', value: amountType(category.amountType)},
            {input: category, key: 'security', value: null, name: 'displays blank for security = false'},
            {input: category, key: 'income', value: null, name: 'displays blank for income = false'},
            {input: category, key: 'transactionCount', value: category.transactionCount},
            {input: newCategoryModel({security: true}), key: 'security', value: <span>&#x1f5f8;</span>, name: 'displays check for security = true'},
            {input: newCategoryModel({income: true}), key: 'income', value: <span>&#x1f5f8;</span>, name: 'displays check for income = true'},
        ];

        columnTests.forEach(({name, input, key, value}) => {
            it(name ?? `displays ${key}`, () => {
                const component = shallow(<CategoriesPage />);
                const columns = component.find(Table).prop<IColumn<CategoryModel>[]>('columns');

                const column = columns.find((column) => column.key === `category.${key}`);

                expect(column.render(input)).toEqual(value);
            });
        });
    });
});
