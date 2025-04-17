import { useEffect, useState } from "react";
import axios from "axios";
import { Table } from "@/components/ui/table";
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer } from "recharts";

export default function AdminPanel() {
    const [data, setData] = useState([]);
    
    useEffect(() => {
        axios.get("http://localhost:8080/analytics")
            .then(response => setData(response.data))
            .catch(error => console.error("Error fetching analytics:", error));
    }, []);

    return (
        <div className="p-6">
            <h1 className="text-xl font-bold mb-4">URL Shortener Analytics</h1>
            
            {/* Table for Analytics Data */}
            <Table>
                <thead>
                    <tr>
                        <th>Short URL</th>
                        <th>Long URL</th>
                        <th>Clicks</th>
                    </tr>
                </thead>
                <tbody>
                    {data.map((url, index) => (
                        <tr key={index}>
                            <td>
                                <a href={`http://localhost:8080/${url.short_url}`} target="_blank" rel="noopener noreferrer">
                                    {url.short_url}
                                </a>
                            </td>
                            <td className="truncate max-w-xs">{url.long_url}</td>
                            <td>{url.clicks}</td>
                        </tr>
                    ))}
                </tbody>
            </Table>

            {/* Clicks Chart */}
            <div className="mt-6">
                <h2 className="text-lg font-bold mb-2">Click Trends</h2>
                <ResponsiveContainer width="100%" height={300}>
                    <BarChart data={data} margin={{ top: 10, right: 30, left: 0, bottom: 10 }}>
                        <XAxis dataKey="short_url" />
                        <YAxis />
                        <Tooltip />
                        <Bar dataKey="clicks" fill="#4F46E5" />
                    </BarChart>
                </ResponsiveContainer>
            </div>
        </div>
    );
}
