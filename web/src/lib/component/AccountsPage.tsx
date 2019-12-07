import React from 'react';
import {Link} from 'react-router-dom';
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
                        <th className='enum'>Type</th>
                        <th>Description</th>
                        <th>Account Number</th>
                        <th className='boolean'>Closed</th>
                        <th className='number'>Transactions</th>
                        <th className='number'>Balance</th>
                    </tr>
                </thead>
                <tbody>
                    {accounts.map(account => (
                        <tr>
                            <td>{account.companyName}</td>
                            <td><Link to={`account/${account.id}`}>{account.name}</Link></td>
                            <td className='enum'>{account.type}</td>
                            <td>{account.description || ''}</td>
                            <td>{account.accountNo || ''}</td>
                            <td className='boolean'>{account.closed ? <span>&#x1F5F8;</span> : null}</td>
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