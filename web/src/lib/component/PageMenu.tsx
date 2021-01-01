import React from 'react';
import {Link} from 'react-router-dom';
import MenuItem from '@material-ui/core/MenuItem';
import MenuList from '@material-ui/core/MenuList';
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

const PageMenu = React.forwardRef<HTMLUListElement, IProps>(({currentPage}, ref) => {
    const items = currentPage ? menuItems.filter(([key]) => currentPage !== key) : menuItems;
    return (
        <MenuList ref={ref}>
            {items.map(([key, url]) => (<MenuItem key={key}><Link to={url}>{translate(key)}</Link></MenuItem>))}
        </MenuList>
    );
});

export default PageMenu;
