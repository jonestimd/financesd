import React from 'react';
import {RootStoreContext} from '../store/RootStore';
import {observer} from 'mobx-react-lite';

const AccountsPage: React.FC<{}> = observer(() => {
    const rootStore = React.useContext(RootStoreContext);
    const accounts = rootStore.accountStore.accounts;
    return (
        <div className='account-list'>
            <table className='table'>
                <thead>
                    <tr>
                        <th>Closed</th>
                        <th>Company</th>
                        <th>Name</th>
                        <th>Type</th>
                        <th>Description</th>
                        <th>Account Number</th>
                    </tr>
                </thead>
                <tbody>
                    {accounts.map(account => (
                        <tr>
                            <td className='check-mark'>{account.closed ? <span>&#x1F5F8;</span> : null}</td>
                            <td>{account.companyName}</td>
                            <td>{account.name}</td>
                            <td>{account.type}</td>
                            <td>{account.description || ''}</td>
                            <td>{account.accountNo || ''}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
});

export default AccountsPage;