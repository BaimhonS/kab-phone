import React, { useState, useCallback } from 'react';
import { Search } from 'lucide-react';
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { debounce } from 'lodash';

const SearchBar = ({ onSearch }) => {
  const [searchTerm, setSearchTerm] = useState('');

  const debouncedSearch = useCallback(
    debounce((value) => {
      onSearch(value);
    }, 1000),
    [onSearch]
  );

  const handleInputChange = (e) => {
    const value = e.target.value;
    setSearchTerm(value);
    debouncedSearch(value);
  };

  return (
    <div className="container mx-auto px-4 py-4">
      <div className="flex gap-2">
        <Input 
          type="text" 
          placeholder="Search here..." 
          className="flex-grow"
          value={searchTerm}
          onChange={handleInputChange}
        />
        <Button variant="outline" size="icon">
          <Search className="h-4 w-4" />
        </Button>
      </div>
    </div>
  );
};

export default SearchBar;