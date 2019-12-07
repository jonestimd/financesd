import React from 'react';
import {RootStoreContext} from '../store/RootStore';
import {observer} from 'mobx-react-lite';
import * as formats from '../formats';

const AccountsPage: React.FC<{}> = observer(() => {
    const {accountStore} = React.useContext(RootStoreContext);
    const accounts = accountStore.accounts;
    React.useEffect(() => accountStore.getAccounts(), []);
    return (
        <div className='account-list'>
            <table className='table'>
                <thead>
                    <tr>
                        <th>Company</th>
                        <th>Name</th>
                        <th>Type</th>
                        <th>Description</th>
                        <th>Account Number</th>
                        <th>Closed</th>
                        <th>Transactions</th>
                        <th>Balance</th>
                    </tr>
                </thead>
                <tbody>
                    {accounts.map(account => (
                        <tr>
                            <td>{account.companyName}</td>
                            <td>{account.name}</td>
                            <td>{account.type}</td>
                            <td>{account.description || ''}</td>
                            <td>{account.accountNo || ''}</td>
                            <td className='check-mark'>{account.closed ? <span>&#x1F5F8;</span> : null}</td>
                            <td className='number'>{account.transactionCount}</td>
                            <td className='number'>{formats.currency.format(account.balance)}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
});

export default AccountsPage;