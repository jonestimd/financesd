import React from 'react';
import {RootStoreContext} from '../store/RootStore';
import {observer} from 'mobx-react-lite';
import TopAppBar from './TopAppBar';
import Table from './table/Table';
import {IColumn} from './table/Column';
import {CategoryModel} from '../model/CategoryModel';
import amountType from '../i18n/amountType';
import {translate} from '../i18n/localize';

const columns: IColumn<CategoryModel>[] = [
    {key: 'category.parent', render: (category) => category.parent?.displayName ?? null},
    {key: 'category.code', render: (category) => category.code},
    {key: 'category.description', render: (category) => category.description},
    {key: 'category.amountType', render: (category) => amountType(category.amountType)},
    {key: 'category.security', render: (category) => category.security ? <span>&#x1f5f8;</span> : null, className: 'boolean'},
    {key: 'category.income', render: (category) => category.income ? <span>&#x1f5f8;</span> : null, className: 'boolean'},
    {key: 'category.transactionCount', render: (category) => category.transactionCount, className: 'number'},
];

const CategoriesPage: React.FC = observer(() => {
    const {categoryStore} = React.useContext(RootStoreContext);
    const categories = categoryStore.categories;
    return <>
        <TopAppBar title={translate('menu.categories')} currentPage='menu.categories' />
        <Table columns={columns} data={categories} />
    </>;
});

export default CategoriesPage;
