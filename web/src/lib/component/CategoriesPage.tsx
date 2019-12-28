import React from 'react';
import {RootStoreContext} from '../store/RootStore';
import {observer} from 'mobx-react-lite';
import TopAppBar from './TopAppBar';
import Table, {IColumn} from './Table';
import {CategoryModel} from '../model/CategoryModel';
import {translate} from '../i18n/localize';

const columns: IColumn<CategoryModel>[] = [
	{key: 'category.parent', render: (category) => category.displayName},
    {key: 'category.code', render: (category) => category.code},
    {key: 'category.description', render: (category) => category.description},
    {key: 'category.amountType', render: (category) => category.amountType},
    {key: 'category.security', render: (category) => category.security ? <span>&#x1f5f8;</span> : null, className: 'boolean'},
    {key: 'category.income', render: (category) => category.income ? <span>&#x1f5f8;</span> : null, className: 'boolean'},
    {key: 'category.transactionCount', render: (category) => category.transactionCount, className: 'number'},
];

const CategoriesPage: React.FC<{}> = observer(() => {
    const menuItems = [translate('menu.accounts'), translate('menu.payees'), translate('menu.securities')];
    const {categoryStore} = React.useContext(RootStoreContext);
    const categories = categoryStore.categories;
    React.useEffect(() => categoryStore.loadCategories(), []);
    return (
        <div className='category-list'>
            <TopAppBar title={translate('menu.categories')} menuItems={menuItems} />
            <Table columns={columns} data={categories} />
        </div>
    );
});

export default CategoriesPage;