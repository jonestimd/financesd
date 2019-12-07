import React from 'react';
import {Link} from 'react-router-dom';
import {RootStoreContext} from '../store/RootStore';
import {observer} from 'mobx-react-lite';
import * as formats from '../formats';
import Table, {IColumn} from './Table';
import {AccountModel} from '../model/AccountModel';

const columns: IColumn<AccountModel>[] = [
    {name: 'Company', getter: (account) => account.companyName},
    {name: 'Name', getter: (account) => <Link to={`account/${account.id}`}>{account.name}</Link>},
    {name: 'Type', getter: (account) => account.type, className: 'enum'},
    {name: 'Description', getter: (account) => account.description},
    {name: 'Account Number', getter: (account) => account.accountNo},
    {name: 'Closed', getter: (account) => account.closed ? <span>&#x1F5F8;</span> : null, className: 'boolean'},
    {name: 'Transactions', getter: (account) => account.transactionCount, className: 'number'},
    {name: 'Balance', getter: (account) => formats.currency.format(account.balance), className: 'number'},
];

const AccountsPage: React.FC<{}> = observer(() => {
    const {accountStore} = React.useContext(RootStoreContext);
    const accounts = accountStore.accounts;
    React.useEffect(() => accountStore.getAccounts(), []);
    return (
        <div className='account-list'>
            <Table columns={columns} data={accounts}/>
        </div>
    );
});

export default AccountsPage;