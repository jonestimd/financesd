import React from 'react';
import {Link} from 'react-router-dom';
import TopAppBar from './TopAppBar';
import {translate} from '../i18n/localize';

interface IProps {
    match: {params: {[name: string]: string}}
}

const TransactionsPage: React.FC<IProps> = ({match: {params: {accountId}}}) => {
    const menuItems = [
        <Link to='/finances/' className='menu-item'>{translate('menu.accounts')}</Link>
    ]
    return (
        <TopAppBar title={'Account #' + accountId} menuItems={menuItems}/>
    );
};

export default TransactionsPage;