import React from 'react';
import {Link} from 'react-router-dom';
import MenuItem from '@material-ui/core/MenuItem';
import {translate} from '../i18n/localize';

const menuItems = [
    ['menu.accounts', '/finances/'],
    ['menu.categories', '/finances/categories'],
    ['menu.payees', '/finances/payees'],
    ['menu.securities', '/finances/securities'],
];

interface IProps {
    currentPage?: string;
}

const PageMenu: React.FC<IProps> = ({currentPage}) => {
    const items = currentPage ? menuItems.filter(([key]) => currentPage !== key) : menuItems;
    return <>{items.map(([key, url]) => (
        <MenuItem key={key}><Link to={url}>{translate(key)}</Link></MenuItem>
    ))}</>;
};

export default PageMenu;
