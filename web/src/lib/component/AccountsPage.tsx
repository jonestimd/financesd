import React from 'react';
import {Link} from 'react-router-dom';
import {RootStoreContext} from '../store/RootStore';
import {observer} from 'mobx-react-lite';
import * as formats from '../formats';
import TopAppBar from './TopAppBar';
import Table, {IColumn} from './Table';
import {AccountModel} from '../model/AccountModel';
import accountType from '../i18n/accountType';
import {translate} from '../i18n/localize';

const columns: IColumn<AccountModel>[] = [
    {name: translate('account.company'), getter: (account) => account.companyName},
    {name: translate('account.name'), getter: (account) => <Link to={`account/${account.id}`}>{account.name}</Link>},
    {name: translate('account.type'), getter: (account) => accountType(account.type), className: 'enum'},
    {name: translate('account.description'), getter: (account) => account.description},
    {name: translate('account.number'), getter: (account) => account.accountNo},
    {name: translate('account.closed'), getter: (account) => account.closed ? <span>&#x1F5F8;</span> : null, className: 'boolean'},
    {name: translate('account.transactions'), getter: (account) => account.transactionCount, className: 'number'},
    {name: translate('account.balance'), getter: (account) => formats.currency.format(account.balance), className: 'number'},
];

const AccountsPage: React.FC<{}> = observer(() => {
    const menuItems = [translate('menu.categories'), translate('menu.securities')];
    const {accountStore} = React.useContext(RootStoreContext);
    const accounts = accountStore.accounts;
    React.useEffect(() => accountStore.getAccounts(), []);
    return (
        <div className='account-list'>
            <TopAppBar title={translate('menu.accounts')} menuItems={menuItems} />
            <Table columns={columns} data={accounts} />
        </div>
    );
});

export default AccountsPage;