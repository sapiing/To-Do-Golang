"use client";

import { useState, useEffect, useCallback } from "react";
import FormTask from "../components/form_task";

export default function TaskPage() {
  const [showModal, setShowModal] = useState(false);
  const [tasks, setTasks] = useState([]); // State to store fetched tasks
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const [editingTask, setEditingTask] = useState(null);

  const openModal = () => setShowModal(true);
  const closeModal = () => {
    setShowModal(false);
    setEditingTask(null);
  };

  //Fetch tasks from backend on component mount

  function handleEdit(task) {
    setEditingTask(task); // Set the task to be edited
    openModal(); // Open the modal for editing
  }

  const fetchTasks = useCallback(async () => {
    try {
      const token = localStorage.getItem("token"); // Get token from local storage
      const response = await fetch("http://localhost:8080/api/task", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error(`Request failed with status ${response.status}`);
    }
    

      const data = await response.json();

      const formattedTasks = data.map((task) => ({
        ...task,
        date: new Date(task.date).toISOString(), // Adjust to your preferred date format
      }));
      setTasks(formattedTasks);

      setTasks(data);
    } catch (error) {
      setError(error.message);
    } finally {
      setIsLoading(false);
    }
  }, []);

  async function handleDelete(taskId) {
    try {
      const token = localStorage.getItem("token");
      const response = await fetch(`http://localhost:8080/api/task/${taskId}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error("Failed to delete task");
      }

      // Update the task list after successful deletion
      fetchTasks();
    } catch (error) {
      setError(error.message);
    }
  }

  async function handleToggleComplete(task) {
    try {
      const token = localStorage.getItem("token");
      const response = await fetch(
        `http://localhost:8080/api/task/${task.id}`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            id: task.id,
            completed: !task.completed,
          }), // Toggle completed
        }
      );

      if (!response.ok) {
        throw new Error("Failed to update task status");
      }

      // Update the task list after successful status change
      fetchTasks();
    } catch (error) {
      setError(error.message);
    }
  }

  useEffect(() => {
    fetchTasks();
  }, [fetchTasks]); // Empty dependency array ensures this runs only once on mount

  return (
    <div>
      <button
        onClick={openModal}
        className="bg-teal-500 hover:bg-teal-600 transition-colors duration-300 text-white font-semibold py-2.5 px-6 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-teal-400"
      >
        Tambah Task
      </button>

      {isLoading ? (
        <p className="text-gray-500 mt-4">Loading tasks...</p>
      ) : error ? (
        <p className="text-red-500 mt-4">Error: {error}</p>
      ) : (
        <table className="table-auto w-full mt-6 border-collapse">
          <thead className="text-left bg-gray-100">
            <tr className="rounded-lg">
              <th className="py-3 px-4 font-semibold rounded-tl-lg">Title</th>
              <th className="py-3 px-4 font-semibold">Description</th>
              <th className="py-3 px-4 font-semibold rounded-tr-lg">Actions</th>
            </tr>
          </thead>
          <tbody>
            {tasks.length === 0 ? (
              <tr>
                <td colSpan="4" className="py-4 text-center text-gray-500">
                  No tasks found. Create a new task to get started!
                </td>
              </tr>
            ) : (
              tasks.map((task) => (
                <tr key={task.id} className="hover:bg-gray-50">
                  <td className="py-2 px-4 border-b">{task.title}</td>

                  <td className="py-2 px-4 border-b">{task.description}</td>
                  <td className="py-2 px-4 border-b space-x-2">
                    <button
                      onClick={() => handleEdit(task)}
                      className="bg-blue-500 hover:bg-blue-600 transition-colors duration-300 text-white font-semibold py-1.5 px-3 rounded-md text-sm"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => handleDelete(task.id)}
                      className="bg-red-500 hover:bg-red-600 transition-colors duration-300 text-white font-semibold py-1.5 px-3 rounded-md text-sm"
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      )}

      {showModal && (
        <div className="fixed z-10 inset-0 overflow-y-auto">
          <div className="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
            {/* Background Overlay */}
            <div
              className="fixed inset-0 transition-opacity"
              aria-hidden="true"
            >
              <div className="absolute inset-0 bg-gray-500 opacity-75"></div>
            </div>

            {/* Modal Content */}
            <span
              className="hidden sm:inline-block sm:align-middle sm:h-screen"
              aria-hidden="true"
            >
              &#8203;
            </span>
            <div className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
              <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
                {/* ... FormTask component is placed here ... */}
                <FormTask
                  closeModal={closeModal}
                  refetchTasks={fetchTasks}
                  editingTask={editingTask}
                />
                {/*Pass the closeModal function as a prop */}
              </div>
              <div className="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
                <button
                  type="button"
                  onClick={closeModal}
                  className="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
                >
                  Cancel
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
