import React from 'react';
import {Link} from 'react-router-dom';
import {MenuListItem, MenuListItemText} from '@material/react-menu';
import {translate} from '../i18n/localize';

const menuItems = [
    ['menu.accounts', '/finances/'],
    ['menu.categories', '/finances/categories'],
    ['menu.payees', '/finances/payees'],
    ['menu.securities', '/finances/securities'],
];

interface IProps {
    currentPage: string;
}

const PageMenu: React.FC<IProps> = ({currentPage}) => {
    const items = menuItems.filter(([key]) => currentPage !== key);
    return <>{items.map(([key, url]) => (
        <MenuListItem key={key}>
            <MenuListItemText primaryText={<Link to={url}>{translate(key)}</Link>} />
        </MenuListItem>
    ))}</>;
};

export default PageMenu;
