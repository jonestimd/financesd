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
                    <tr><th>Company</th><th>Name</th></tr>
                </thead>
                <tbody>
                    {accounts.map(account => (
                        <tr>
                            <td>{account.companyName}</td>
                            <td>{account.name}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
});

export default AccountsPage;