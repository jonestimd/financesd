import React from 'react';
import {Link} from 'react-router-dom';
import {RootStoreContext} from '../store/RootStore';
import {observer} from 'mobx-react-lite';
import * as formats from '../formats';
import TopAppBar from './TopAppBar';
import Table from './table/Table';
import {IColumn} from './table/Column';
import {AccountModel} from '../model/account/AccountModel';
import accountType from '../i18n/accountType';
import {translate} from '../i18n/localize';
import CompaniesDialog from './CompaniesDialog';
import {Icon, IconButton} from '@material-ui/core';

const columns: IColumn<AccountModel>[] = [
    {key: 'account.company', render: (account) => account.companyName},
    {key: 'account.name', render: (account) => <Link to={`account/${account.id}`}>{account.name}</Link>},
    {key: 'account.type', render: (account) => accountType(account.type), className: 'enum'},
    {key: 'account.description', render: (account) => account.description},
    {key: 'account.number', render: (account) => account.accountNo},
    {key: 'account.closed', render: (account) => account.closed ? <span>&#x2713;</span> : null, className: 'boolean'},
    {key: 'account.transactions', render: (account) => account.transactionCount, className: 'number'},
    {key: 'account.balance', render: (account) => formats.currency.format(account.balance), className: 'number'},
];

const AccountsPage: React.FC = observer(() => {
    const {accountStore} = React.useContext(RootStoreContext);
    const accounts = accountStore.accounts;
    const [showCompanies, setShowCompanies] = React.useState(false);
    return (
        <div className='account-list'>
            <TopAppBar title={translate('menu.accounts')} currentPage='menu.accounts'>
                <IconButton onClick={() => setShowCompanies(true)}><Icon>account_balance</Icon></IconButton>
            </TopAppBar>
            {showCompanies ? <CompaniesDialog onClose={() => setShowCompanies(false)} /> : null}
            <Table columns={columns} data={accounts} />
        </div>
    );
}, {forwardRef: true});

export default AccountsPage;
