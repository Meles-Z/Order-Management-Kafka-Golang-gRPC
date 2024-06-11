import '../../assets/styles/header.css'
export default function Header() {
    return (
        <>
            <div className="all-header bg-red-700">
                <div className="container mx-auto flex  items-center py-4">
                    <div className="logo flex-shrink-0">
                        <img src="" alt="image is left here" className="h-8 text-white" />
                    </div>
                    <div className="flex-grow flex justify-center">
                        <ul className="flex gap-8 text-white text-lg">
                            <li>Popular</li>
                            <li>Cars</li>
                            <li>Sellers</li>
                            <li>Sell</li>
                            <li>Search</li>
                        </ul>
                    </div>
                    <div className="flex-shrink-0">
                        <ul className="flex gap-4 items-center text-white items-center ... ">
                            <li >
                                <button>
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
                                <path stroke-linecap="round" stroke-linejoin="round" d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z" />
                                </svg>
                                </button>
                               </li>
                            <li >
                                <button>
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M17.982 18.725A7.488 7.488 0 0 0 12 15.75a7.488 7.488 0 0 0-5.982 2.975m11.963 0a9 9 0 1 0-11.963 0m11.963 0A8.966 8.966 0 0 1 12 21a8.966 8.966 0 0 1-5.982-2.275M15 9.75a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
                                </svg>
                                </button>
                            </li>
                            <li>
                                <button className="flex gap-2 btn-add-car">
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
                                <path strokeLinecap="round" strokeLinejoin="round" d="M12 9v6m3-3H9m12 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
                                </svg>
                                    Add Car
                                </button>
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
        </>
    );
}
